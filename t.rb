require "uri"
require "net/http"

url = URI("http://localhost:8080/servers")
http = Net::HTTP.new(url.host, url.port)

request = Net::HTTP::Post.new(url)
request["Content-Type"] = "application/x-www-form-urlencoded"
request["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiaXNzIjoia2stb3BzIiwic3ViIjoiYWRtaW4iLCJleHAiOjE3NDc5MjUzNjQsIm5iZiI6MTc0NzgzODk2NCwiaWF0IjoxNzQ3ODM4OTY0fQ.D5yLdigkGwY07Rnn8idKv71CTfkQmuBTB7Hmmh_1uv0"  # 可选：从 Cookie 中提取的 token

# 请求体参数
form_data = {
  name: "node04",
  ip: "1.1.1.4",
  port: 22,
  username: "kk",
  status: "online",
  auth_type: "key",
  password: "",
  ssh_key_path: "~/ssh/id_rsa",
  description: ""
}

# 编码表单数据
request.body = URI.encode_www_form(form_data)

response = http.request(request)

puts "Response Code: #{response.code}"
puts "Response Body: #{response.body}"
