package main


import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/google/jsonapi"
	"net/http"
	"time"
	"api_test_2/db"
	"encoding/json"
	"api_test_2/utils"
	"api_test_2/models"
	"github.com/urfave/negroni"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	//"log"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)


type CategoryJSON struct {
	Name		string
	ParentID	uint
}

type ProductJSON struct {
	Name		string
	Price		float32
}

var SECRET_TOKEN = "LongSecretString"  // Secret string for JWT

func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Index!")
}

var NewCategoryHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		user := req.Context().Value("user");
		claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
		userID_float := claims["user_id"].(float64)
		userID := int(userID_float)
		IsAdmin := utils.UserIsAdmin(userID)

		if IsAdmin {
			decoder := json.NewDecoder(req.Body)

			var categoryJson CategoryJSON
			err := decoder.Decode(&categoryJson)
			if err != nil {
				panic(err)
				return
			}

			utils.AddCategory(categoryJson.Name, categoryJson.ParentID)
		} else {
			var err_map = map[string]string{"error": "Access denied. Need admin access."}
			json.NewEncoder(w).Encode(err_map)
		}
	}
})

var DeleteCategoryHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user");
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))
	IsAdmin := utils.UserIsAdmin(userID)

	if IsAdmin {
		categoryID, err := strconv.Atoi(req.URL.Query().Get("category_id"))
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Wrong category_id. Please check.."})
			return
		}

		utils.DeleteCategory(categoryID)
	} else {
		var err_map = map[string]string{"error": "Access denied. Need admin access."}
		json.NewEncoder(w).Encode(err_map)
	}
})

var CategoriesHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	databaseConnect := db.GetDatabaseConnectInstance()
	dbConnect := databaseConnect.GetDbConnect()

	jsonapiRuntime := jsonapi.NewRuntime().Instrument("categories.list")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", jsonapi.MediaType)

	var categories []*models.Category
	dbConnect.Find(&categories)

	if err := jsonapiRuntime.MarshalPayload(w, categories); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
})

var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var filterCategoryID int
	filterCategoryIDRaw := req.URL.Query().Get("category_id")
	if filterCategoryIDRaw != "" {
		filterCategoryID, _ = strconv.Atoi(filterCategoryIDRaw)
	}

	fmt.Println("Products")
	databaseConnect := db.GetDatabaseConnectInstance()
	dbConnect := databaseConnect.GetDbConnect()

	jsonapiRuntime := jsonapi.NewRuntime().Instrument("products.list")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", jsonapi.MediaType)

	var products []*models.Product
	if filterCategoryID != 0 {
		// Filter by category id
		dbConnect.Where("category_id = ?", filterCategoryID).Find(&products)
	} else {
		dbConnect.Find(&products)
	}

	if err := jsonapiRuntime.MarshalPayload(w, products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
})

var ProductHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Product")

	jsonapiRuntime := jsonapi.NewRuntime().Instrument("products.detail")

	vars := mux.Vars(req)
	id := vars["id"]
	var product models.Product
	databaseConnect := db.GetDatabaseConnectInstance()
	dbConnect := databaseConnect.GetDbConnect()

	dbConnect.Find(&product, id)

	if err := jsonapiRuntime.MarshalPayload(w, &product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
})

var UpdateProductHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		user := req.Context().Value("user");
		claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
		userID_float := claims["user_id"].(float64)
		userID := int(userID_float)
		IsAdmin := utils.UserIsAdmin(userID)

		if IsAdmin {
			decoder := json.NewDecoder(req.Body)

			var productJSON ProductJSON
			err := decoder.Decode(&productJSON)
			if err != nil {
				panic(err)  // Not need panic
				return
			}

			vars := mux.Vars(req)
			productID, err := strconv.Atoi(vars["id"])
			if err != nil {
				fmt.Println(err)  // Can't cast string id to int?
				return
			}

			//fmt.Println(productJSON.Name)
			//fmt.Println(productJSON.Price)
			utils.UpdateProduct(productID, productJSON.Name, productJSON.Price)
		} else {
			var err_map = map[string]string{"error": "Access denied. Need admin access."}
			json.NewEncoder(w).Encode(err_map)
		}
	}
})

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		login := req.FormValue("login")
		//fmt.Println(login)

		password := req.FormValue("password")

		if login != "" && password != "" {

			var err_map = map[string]string{"error": "Wrong username or password"}

			var user models.User

			databaseConnect := db.GetDatabaseConnectInstance()
			dbConnect := databaseConnect.GetDbConnect()

			result_user := dbConnect.Where("login = ?", login).First(&user)
			//fmt.Println(user)
			if result_user.Error != nil {
				json.NewEncoder(w).Encode(err_map)
				fmt.Println(result_user.Error)
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
				fmt.Println(err)
				json.NewEncoder(w).Encode(err_map)
				return
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user_id": user.ID,
				"admin": user.IsAdmin,
			})

			tokenString, err := token.SignedString([]byte(SECRET_TOKEN))

			if err != nil {
				return  // Need send error?
			}

			tokenResponse := map[string]string{"token": tokenString}

			json.NewEncoder(w).Encode(tokenResponse)
			//fmt.Println("All right!")
		}
		//fmt.Println(password)
	}
}

func initDatabase() {
	// Simple fixtures (initial data for database)
	fixtures_category := [6]models.Category{
		models.Category{Name:"Computers", Lft: 0, Rgt: 5, CreatedAt: time.Now(), UpdatedAt: time.Now(), Level: 0, ParentID: 0},
		models.Category{Name:"Notebooks", Lft: 1, Rgt: 2, CreatedAt: time.Now(), UpdatedAt: time.Now(), Level: 1, ParentID: 1},
		models.Category{Name:"Tablets", Lft: 3, Rgt: 4, CreatedAt: time.Now(), UpdatedAt: time.Now(), Level: 1, ParentID: 1},
		models.Category{Name:"Monitors", Lft: 6, Rgt: 11, CreatedAt: time.Now(), UpdatedAt: time.Now(), Level: 0, ParentID: 0},
		models.Category{Name:"TN", Lft: 7, Rgt: 8, CreatedAt: time.Now(), UpdatedAt: time.Now(), Level: 1, ParentID: 3},
		models.Category{Name:"IPS", Lft: 9, Rgt: 10, CreatedAt: time.Now(), UpdatedAt: time.Now(), Level: 1, ParentID: 3},
	}

	fixtures_product := [3]models.Product{
		models.Product{Name:"ASUS Notebook 1", CategoryID: 2, Price: 30000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Product{Name:"ASUS Notebook 2", CategoryID: 2, Price: 25000.5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Product{Name:"Acer Tablet 1", CategoryID: 3, Price: 10000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	fixtures_users := [2]models.User{
		models.User{Login: "admin", Password: "$2a$10$thUwUVkEQ.ewvTSCdQ.38eoKaSAkt0plfwk9ufPL74YKJLbN7urjG", IsAdmin: true},
		models.User{Login: "user", Password: "$2a$10$CoK1nVAL2t3nifV8M8RoVO6nwxMeo/p6nOVC9yCTuD3nPxkl6lsDO", IsAdmin: false},
	}

	database := db.NewDbConnect("root:qwqw4@tcp(localhost:3306)/api_test_2?parseTime=true")
	databaseConnect := database.GetDbConnect()

	// Migrate models
	databaseConnect.AutoMigrate(&models.Category{})
	databaseConnect.AutoMigrate(&models.Product{})
	databaseConnect.AutoMigrate(&models.User{})

	// Get count of all categories
	var cntCategories int
	result_cnt_categories := databaseConnect.Model(&models.Category{}).Count(&cntCategories)
	if result_cnt_categories.Error != nil {
		fmt.Println(result_cnt_categories.Error)
		return
	}

	// Maybe need apply fixtures? For init database
	if cntCategories == 0 {
		for _, category := range fixtures_category {
			databaseConnect.NewRecord(category)
			databaseConnect.Create(&category)
		}
	}

	var cntProducts int
	result_cnt_products := databaseConnect.Model(&models.Product{}).Count(&cntProducts)
	if result_cnt_products.Error != nil {
		fmt.Println(result_cnt_products.Error)
		return
	}
	if cntProducts == 0 {
		for _, product := range fixtures_product {
			databaseConnect.NewRecord(product)
			databaseConnect.Create(&product)
		}
	}

	var cntUsers int
	result_cnt_users := databaseConnect.Model(&models.User{}).Count(&cntUsers)
	if result_cnt_users.Error != nil {
		fmt.Println(result_cnt_users.Error)
		return
	}

	if cntUsers == 0 {
		for _, user := range fixtures_users {
			databaseConnect.NewRecord(user)
			databaseConnect.Create(&user)
		}
	}
}

func main() {
	// Initialize database (migrations, fixtures)
	initDatabase()

	// JWT Middleware for check JWT token (correct or not)
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_TOKEN), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	router := mux.NewRouter()
	router.HandleFunc("/", Index)
	router.HandleFunc("/login/", LoginHandler)

	// Route to create new category (only for admin)
	router.Handle("/new_category/", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(NewCategoryHandler),
	))

	// Route for delete category
	router.Handle("/delete_category/", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(DeleteCategoryHandler),
	))

	// Route for get all products
	router.Handle("/products/", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(ProductsHandler),
	))

	// Route for get product by ID
	router.Handle("/products/{id:[0-9]+}/", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(ProductHandler),
	))

	// Route for update product by ID
	router.Handle("/update_product/{id:[0-9]+}/", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(UpdateProductHandler),
	))

	// Route for get all categories
	router.Handle("/categories/", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(CategoriesHandler),
	))

	http.ListenAndServe(":8080", router)
}