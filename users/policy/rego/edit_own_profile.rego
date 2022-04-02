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
    "deleted_at" not in input.subject
}

# Can edit own username if:
# - can modify their profile
# - never changed their username
edit_username = true {
    allowed
	count(input.subject.update_history.usernames) == 0
}

# Can edit own username if:
# - can modify their profile
# - last change was made more than x days ago
# TODO: how to pass the days as a variable? (maybe from settings and always constructing the variable from file)
edit_username = true {
    allowed
    last_change := time.parse_rfc3339_ns(input.subject.update_history.usernames[0].time)
	(time.now_ns() - last_change) / (60*60*24*1000000000) > input.profile_settings.username_update_days
}