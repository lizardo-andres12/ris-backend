import io
import requests
from PIL import Image

def compress_img(path, quality=90):
    with Image.open(path) as img:
        # Convert to RGB if necessary (e.g. if source is PNG/RGBA) to save as JPEG
        if img.mode in ("RGBA", "P"):
            img = img.convert("RGB")
            
        image_bytes = io.BytesIO()
        img.save(image_bytes, format='JPEG', quality=quality, optimize=True)
        image_bytes.seek(0)
        return image_bytes

def send_bytes(bytes_arr, dest_url):
    try:
        resp = requests.post(
            dest_url, 
            files={'image': bytes_arr}
        )
        
        if resp.status_code == 200:
            print(f'Success from client: {resp.text}')
        else:
            print(f'Error from server [{resp.status_code}]: {resp.text}')
    except Exception as e:
        print(f'An error occurred: {e}')

if __name__ == '__main__':
    img_data = compress_img('./apple.jpg')
    send_bytes(img_data, 'http://localhost:65000')

