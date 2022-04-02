package users

default allowed = false
default edit_username = false

# Actor: logged user
# Subject: person to modify

# User can modify a profile if:
# - They are logged in
# - The person they want to modify has the same id
# TODO: other situations
allowed = true {
	not input.subject.is_banned
    input.actor.user_id == input.subject.user_id
}