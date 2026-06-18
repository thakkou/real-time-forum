{
    "status_code": 200,
    "message": "post deleted successfully",
    "data": {
        "deleted": true,
        "post_id": 1
    }
}
for delete /api/posts/{id}/delete


for creating succes /api/posts/create
what i should get 
{
  "title": "My First Go Post",
  "text": "This is a test post created from Postman using JSON.",
  "categories": [
    "General",
    "Lifestyle",
    "Education"
  ],
image:"https://kslsls" //optional
}

whart i return if suces
{
    "status_code": 200,
    "message": "post created successfully",
    "data": {
        "categories": [
            "General",
            "Lifestyle",
            "Education"
        ],
        "post_id": 4,
        "text": "This is a test post created from Postman using JSON.",
        "title": "My First Go Post",
        "user_id": 1
    }
}