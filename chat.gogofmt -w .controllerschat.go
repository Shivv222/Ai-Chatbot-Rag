warning: in the working copy of 'repository/chat_repository.go', LF will be replaced by CRLF the next time Git touches it
[1mdiff --git a/repository/chat_repository.go b/repository/chat_repository.go[m
[1mindex 790c3fb..7da7468 100644[m
[1m--- a/repository/chat_repository.go[m
[1m+++ b/repository/chat_repository.go[m
[36m@@ -62,25 +62,34 @@[m [mfunc GetChatHistory(userID int) ([]models.ChatHistory, error) {[m
 	return chats, nil[m
 }[m
 [m
[31m-func GetRecentChats(userID int, limit int) ([]models.ChatHistory, error) {[m
[32m+[m[32mfunc GetRecentChats(userID int, sessionID int, limit int) ([]models.ChatHistory, error) {[m
 [m
 	query := `[m
 	SELECT id, user_id, user_message, ai_response, created_at[m
 	FROM chat_history[m
 	WHERE user_id = $1[m
[32m+[m	[32mAND session_id = $2[m
 	ORDER BY created_at DESC[m
[31m-	LIMIT $2[m
[32m+[m	[32mLIMIT $3[m
 	`[m
 [m
[31m-	rows, err := config.DB.Query(query, userID, limit)[m
[32m+[m	[32mrows, err := config.DB.Query([m
[32m+[m		[32mquery,[m
[32m+[m		[32muserID,[m
[32m+[m		[32msessionID,[m
[32m+[m		[32mlimit,[m
[32m+[m	[32m)[m
[32m+[m
 	if err != nil {[m
 		return nil, err[m
 	}[m
[32m+[m
 	defer rows.Close()[m
 [m
 	var chats []models.ChatHistory[m
 [m
 	for rows.Next() {[m
[32m+[m
 		var chat models.ChatHistory[m
 [m
 		err := rows.Scan([m
[36m@@ -98,5 +107,9 @@[m [mfunc GetRecentChats(userID int, limit int) ([]models.ChatHistory, error) {[m
 		chats = append(chats, chat)[m
 	}[m
 [m
[32m+[m	[32mif err := rows.Err(); err != nil {[m
[32m+[m		[32mreturn nil, err[m
[32m+[m	[32m}[m
[32m+[m
 	return chats, nil[m
[31m-}[m
\ No newline at end of file[m
[32m+[m[32m}[m
