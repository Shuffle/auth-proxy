
import requests

proxy_config = {
    "http": "http://:pass@localhost:5004/proxy",
    "https": "https://localhost:5004/proxy"
}

url = "http://ipinfo.io/8.8.8.8"
response = requests.get(
        url, 
        proxies=proxy_config,
)

print(response.text)
print("HEADERS: ", response.headers)
print(response.status_code)
