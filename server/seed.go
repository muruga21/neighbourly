package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"neighbourly/server/internal/models"
)

func main() {
	uri := "mongodb+srv://murugaperumalr2004:murugaperumal@muruga.35stdjd.mongodb.net/?retryWrites=true&w=majority&appName=muruga"
	log.Println("Connecting to MongoDB...")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Connection attempt failed.")
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	log.Println("Pinging MongoDB...")
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Ping failed: ", err)
	}
	log.Println("Connected successfully.")

	collection := client.Database("neighbourly").Collection("users")

	providers := []models.User{
		{
			FullName:          "Marcus Sterling",
			Email:             "marcus.sterling@example.com",
			PhoneNumber:       "9876543210",
			Password:          "$2a$10$7Z/U7Pz7u.3z1Z.3z1Z.3z1Z.3z1Z.3z1Z.3z1z1z1z1z1z1z1",
			Role:              "Service Provider",
			IsProfileComplete: true,
			Avatar:            "https://randomuser.me/api/portraits/men/1.jpg",
			Title:             "Master Electrician",
			Experience:        "12 Years",
			BasePrice:         "85",
			Status:            "Available",
			Location: &models.Location{
				Lat: 11.0183, // Gandhipuram
				Lng: 76.9686,
			},
			Services: []models.Service{
				{ID: "s1", Icon: "electrical-services", Title: "Fixture Installation", Desc: "Ceiling fans and LED tracks", Price: "$85"},
				{ID: "s2", Icon: "router", Title: "Smart Home Setup", Desc: "IoT automation and networking", Price: "$120"},
			},
		},
		{
			FullName:          "Sarah Jenkins",
			Email:             "sarah.jenkins@example.com",
			PhoneNumber:       "9876543211",
			Password:          "$2a$10$7Z/U7Pz7u.3z1Z.3z1Z.3z1Z.3z1Z.3z1Z.3z1z1z1z1z1z1z1",
			Role:              "Service Provider",
			IsProfileComplete: true,
			Avatar:            "https://randomuser.me/api/portraits/women/2.jpg",
			Title:             "Professional Cleaner",
			Experience:        "5 Years",
			BasePrice:         "40",
			Status:            "Available",
			Location: &models.Location{
				Lat: 11.0124, // RS Puram
				Lng: 76.9458,
			},
			Services: []models.Service{
				{ID: "s3", Icon: "cleaning-services", Title: "Deep Cleaning", Desc: "Full house sanitation", Price: "$150"},
				{ID: "s4", Icon: "cleaning-services", Title: "Window Cleaning", Desc: "Interior & Exterior glass", Price: "$45"},
			},
		},
		{
			FullName:          "Robert Muller",
			Email:             "robert.muller@example.com",
			PhoneNumber:       "9876543212",
			Password:          "$2a$10$7Z/U7Pz7u.3z1Z.3z1Z.3z1Z.3z1Z.3z1Z.3z1z1z1z1z1z1z1",
			Role:              "Service Provider",
			IsProfileComplete: true,
			Avatar:            "https://randomuser.me/api/portraits/men/3.jpg",
			Title:             "Master Plumber",
			Experience:        "8 Years",
			BasePrice:         "60",
			Status:            "Available",
			Location: &models.Location{
				Lat: 11.0270, // Peelamedu
				Lng: 77.0010,
			},
			Services: []models.Service{
				{ID: "s5", Icon: "plumbing", Title: "Pipe Repair", Desc: "Leak detection and fixing", Price: "$95"},
				{ID: "s6", Icon: "plumbing", Title: "Bathroom Fitting", Desc: "Taps and shower install", Price: "$110"},
			},
		},
		{
			FullName:          "Alisha Khan",
			Email:             "alisha.khan@example.com",
			PhoneNumber:       "9876543213",
			Password:          "$2a$10$7Z/U7Pz7u.3z1Z.3z1Z.3z1Z.3z1Z.3z1Z.3z1z1z1z1z1z1z1",
			Role:              "Service Provider",
			IsProfileComplete: true,
			Avatar:            "https://randomuser.me/api/portraits/women/4.jpg",
			Title:             "Expert Painter",
			Experience:        "10 Years",
			BasePrice:         "75",
			Status:            "Available",
			Location: &models.Location{
				Lat: 10.9850, // Kurichi
				Lng: 76.9630,
			},
			Services: []models.Service{
				{ID: "s7", Icon: "format-paint", Title: "Interior Painting", Desc: "Premium emulsion coat", Price: "$200"},
				{ID: "s8", Icon: "format-paint", Title: "Textured Walls", Desc: "Custom artistic patterns", Price: "$350"},
			},
		},
		{
			FullName:          "David Raj",
			Email:             "david.raj@example.com",
			PhoneNumber:       "9876543214",
			Password:          "$2a$10$7Z/U7Pz7u.3z1Z.3z1Z.3z1Z.3z1Z.3z1Z.3z1z1z1z1z1z1z1",
			Role:              "Service Provider",
			IsProfileComplete: true,
			Avatar:            "https://randomuser.me/api/portraits/men/5.jpg",
			Title:             "Appliance Specialist",
			Experience:        "6 Years",
			BasePrice:         "50",
			Status:            "Available",
			Location: &models.Location{
				Lat: 11.0117, // Singanallur
				Lng: 77.0218,
			},
			Services: []models.Service{
				{ID: "s9", Icon: "build", Title: "AC Service", Desc: "Filter cleaning & gas refill", Price: "$75"},
				{ID: "s10", Icon: "build", Title: "Washing Machine Repair", Desc: "Motor & drum issues", Price: "$90"},
			},
		},
		{
			FullName:          "Priya Sharma",
			Email:             "priya.sharma@example.com",
			PhoneNumber:       "9876543215",
			Password:          "$2a$10$7Z/U7Pz7u.3z1Z.3z1Z.3z1Z.3z1Z.3z1Z.3z1z1z1z1z1z1z1",
			Role:              "Service Provider",
			IsProfileComplete: true,
			Avatar:            "https://randomuser.me/api/portraits/women/6.jpg",
			Title:             "Gardening Expert",
			Experience:        "4 Years",
			BasePrice:         "35",
			Status:            "Available",
			Location: &models.Location{
				Lat: 11.0210, // Vadavalli
				Lng: 76.9015,
			},
			Services: []models.Service{
				{ID: "s11", Icon: "yard", Title: "Lawn Mowing", Desc: "Grass cutting & disposal", Price: "$40"},
				{ID: "s12", Icon: "yard", Title: "Plant Nutrition", Desc: "Organic fertilizing", Price: "$65"},
			},
		},
		{
			FullName:          "Vikram Singh",
			Email:             "vikram.singh@example.com",
			PhoneNumber:       "9876543216",
			Password:          "$2a$10$7Z/U7Pz7u.3z1Z.3z1Z.3z1Z.3z1Z.3z1Z.3z1z1z1z1z1z1z1",
			Role:              "Service Provider",
			IsProfileComplete: true,
			Avatar:            "https://randomuser.me/api/portraits/men/7.jpg",
			Title:             "Professional Carpenter",
			Experience:        "15 Years",
			BasePrice:         "70",
			Status:            "Available",
			Location: &models.Location{
				Lat: 11.0494, // Saravanampatti
				Lng: 77.0194,
			},
			Services: []models.Service{
				{ID: "s13", Icon: "carpenter", Title: "Furniture Repair", Desc: "Tables, chairs, and cabinets", Price: "$80"},
				{ID: "s14", Icon: "carpenter", Title: "Custom Cupboards", Desc: "Modern kitchen woodwork", Price: "$500"},
			},
		},
	}

	for _, p := range providers {
		_, err := collection.InsertOne(context.TODO(), p)
		if err != nil {
			log.Printf("Failed to insert %s: %v", p.FullName, err)
		} else {
			log.Printf("Successfully seeded %s", p.FullName)
		}
	}
}
