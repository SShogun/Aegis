from predict import predict_event


event = {
    "hour": 3,
    "failed_logins": 18,
    "requests_per_minute": 700,
    "files_downloaded": 1200,
}


result = predict_event(event)


print("ML Assessment:")
print()

print(result)