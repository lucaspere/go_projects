new_recipe:
	curl -X POST \
	-H "Content-Type: application/json" \
	-d @./mocks/newRecipe.json \
	http://localhost:8080/recipes \
	| jq -r

list_recipe:
	curl -sL \
	http://localhost:8080/recipes \
	| jq -r