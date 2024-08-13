# Auth Proxy
Auth Proxy is an HTTP proxy that automatically adds the Authorization header to the request. It is useful for testing APIs that require authentication.

## TODO
**Auth proxy is in internal tests by the Shuffle team** and is currently not useful. We will eventually release it as a part of shuffle with the /api/v1/proxy endpoint 

## Authentication Storage
Authentication in Shuffle is handled in many ways, including OAuth2, JWT, and API keys. Auth Proxy supports OAuth2, API keys and every other authentication method Shuffle itself supports.

Shuffle encrypts and stores your authentication data securely. The proxy will automatically add the correct headers to the request, so you don't have to worry about it. 

## Usage
1. Set up auth for your product: shuffler.io/search
2. Get your API-key: shuffler.io/settings 
3. Configure HTTP/HTTPS proxy: 
```bash
export HTTP_PROXY=http://:APIKEY@shuffler.io/api/v1/proxy
export HTTPS_PROXY=https://:APIKEY@shuffler.io/api/v1/proxy
```

4. Run your python script (example): The proxy will automatically add the Authorization header to the request for Github Oauth2 in the case below, as long as your Github authentication in Shuffle is configured. 
```python
import requests

# Manual control of proxies without HTTP/HTTPS_PROXY env variables:
proxies = {"http": "http://:APIKEY@shuffler.io/api/v1/proxy"}

# Send the request
response = requests.get('https://api.github.com/user', proxies=proxies)
print(response.json())
```

## Accessing Local Services 
Since the proxy is in the [Shuffle backend](https://github.com/shuffle/shuffle), you will need a local version of the Shuffle backend running to access local services.

You can redirect requests to the proxy with the `HTTP_PROXY` and `HTTPS_PROXY` environment variables as per usual - just change the `shuffler.io` in the original example to `localhost:5001` or whatever your local Shuffle backend is running on. 
