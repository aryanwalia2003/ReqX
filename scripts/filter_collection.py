import json
import sys

def main():
    file_path = "vuc-dynamic-collection.json"
    with open(file_path, "r", encoding="utf-8") as f:
        data = json.load(f)

    allowed_names = [
        "1. Auth: Receptionist Login",
        "5. Receptionist Selects Clinic",
        "13. Receptionist Broadcast (HTTP)"
    ]

    filtered_requests = []
    for req in data.get("requests", []):
        if req.get("name") in allowed_names:
            filtered_requests.append(req)

    data["requests"] = filtered_requests
    data["name"] = "VUC Dynamic Broadcast Flow"

    with open(file_path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=4)

    print("Successfully filtered collection.")

if __name__ == "__main__":
    main()
