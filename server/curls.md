Upload a file

```bash
curl -X POST http://localhost:3000/upload \
  -F "file=@test.txt" \
  -F "user=asd" \
  -F "location=images"
```
