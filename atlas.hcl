data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path",
    "./models",      # Points to your models directory
    "--dialect",
    "postgres"       # mysql | sqlite | sqlserver
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  
  # The temporary blank database for calculations
  dev = "postgresql://postgres:123@localhost:5432/genzone_dev?sslmode=disable" 
  
  # Your clean main database where the data goes
  url = "postgresql://postgres:123@localhost:5432/genzone?sslmode=disable"

  migration {
    dir = "file://migrations"      
  }
}
