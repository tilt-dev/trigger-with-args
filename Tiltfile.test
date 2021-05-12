# use the new cli to change the args passed to this Cmd! Try:
#   go run cmd/main.go change-me ✨
local_resource('change-me', serve_cmd=['./my-script.sh', '🤖'])

# the new cli can't operate on this Cmd b/c it's wrapped as bash rather than passed as argv!
local_resource('string-cmd', serve_cmd='./my-script.sh 🤐')
