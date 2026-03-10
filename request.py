import requests
import os



current_dir = os.getcwd()
file_path = os.path.join(current_dir, 'task-AFB', 'clients.csv')
url = 'http://localhost:8080/upload'

with open(file_path, 'rb') as f:
    files = {'file': f}
    response = requests.post(url, files=files)