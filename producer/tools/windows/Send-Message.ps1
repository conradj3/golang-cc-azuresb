$body = @{     count = 10 } | ConvertTo-Json  
Invoke-WebRequest -Uri 'http://localhost:8080/createMessages' -Method 'POST' -Body $body -ContentType 'application/json'