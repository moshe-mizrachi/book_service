package constants

var IndexMappings = []struct {
	IndexName string
	Mapping   string
}{
	{
		IndexName: "books",
		Mapping: `
		{
		  "mappings": {
		    "properties": {
		      "title": {
		        "type": "text"
		      },
		      "author_name": {
		        "type": "text"
		      },
		      "price": {
		        "type": "float"
		      },
		      "ebook_available": {
		        "type": "boolean",
		        "null_value": false
		      },
		      "publish_date": {
		        "type": "date",
		        "format": "yyyy-MM-dd"
		      }
		    }
		  }
		}`,
	},
}