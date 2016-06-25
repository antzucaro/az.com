deploy: gen
	rsync --size-only -av public/ azucaro@antzucaro.com:/home/azucaro/antzucaro.com/public/

gen:
	hugo

server:
	hugo server -w

combined-css:
	cat static/static/css/normalize.css static/static/css/skeleton.css static/static/css/blog.css > static/static/css/combined.css

css: combined-css
	yuicompressor --type css -o static/static/css/blog.min.css static/static/css/combined.css
	rm static/static/css/combined.css
