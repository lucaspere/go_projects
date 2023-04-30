get_repo_name: printf '%s\n' "${PWD##*/}"

serve_dev:
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
	MONGO_DATABASE=demo \
	JWT_SECRET=secret \
	go run .

new_recipe:
	curl -X POST \
	-H "Content-Type: application/json" \
	-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNjgyODY4NzQ2fQ.-cyvvlP59yL5MkxiHfmtbXRtVFAPixe5wTTQMR-Dcj0" \
	-d @./mocks/newRecipe.json \
	http://localhost:8080/recipes \
	| jq -r

list_recipe:
	curl -sL \
	http://localhost:8080/recipes \
	| jq -r

signin:
	curl -X POST \
	-H "Content-Type: application/json" \
	-d @./mocks/signIn.json \
	http://localhost:8080/signin \
	| jq -r

signup:
	curl -X POST \
	-H "Content-Type: application/json" \
	-d @./mocks/signup.json \
	http://localhost:8080/signup \
	| jq -r

commit_feature:
	echo "feat($(call list_recipe))"
