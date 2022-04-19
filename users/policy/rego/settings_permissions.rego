package users.settings.permissions

import future.keywords.in

# default allowed = false
# default change_username = false
# default change_dob = false

# allowed = true {
#     not reasons_not_allowed["user_is_banned"]
#     not reasons_not_allowed["user_is_deleted"]
# }

# change_username = true {
#     allowed
#     not input.user.update_history.usernames
# }

# change_username = true {
#     allowed
#     count(input.user.update_history.usernames) == 0
# }

# change_username = true {
#     allowed
#     not reasons_cant_change_username["changed_shortly"]
# }

# change_dob = true {
#     allowed
#     not input.user.update_history.dobs
# }

# change_dob = true {
#     allowed
#     count(input.user.update_history.dobs) == 0
# }

# change_dob = true {
#     allowed
#     not reasons_cant_change_dob["changed_shortly"]
# }

# Reasons why user can't perform an action

not_allowed_change_settings["user_is_banned"] {
    input.user.is_banned
}

not_allowed_change_settings["user_is_deleted"] {
    input.user.deleted_at
}

cant_change_username["not_allowed_change_settings"] {
    count(not_allowed_change_settings) > 0
}

cant_change_username["changed_shortly"] {
    count(not_allowed_change_settings) == 0
    last_change := time.parse_rfc3339_ns(input.user.update_history.usernames[0].time)
	not changed_shortly([last_change, input.profile_settings.username_update_days])
}

cant_change_dob["not_allowed_change_settings"] {
    count(not_allowed_change_settings) > 0
}

cant_change_dob["changed_shortly"] {
    count(not_allowed_change_settings) == 0
    last_change := time.parse_rfc3339_ns(input.user.update_history.dobs[0].time)
	not changed_shortly([last_change, input.profile_settings.dob_update_days])
}

changed_shortly([last_changed, threshold]) = true {
    (time.now_ns() - last_changed) / (60*60*24*1000000000) > threshold
}


# Other stuff I can't think of right now