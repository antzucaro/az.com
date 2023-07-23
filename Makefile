deploy: gen
	rsync --size-only -av public/ ant@antzucaro.com:/var/www/antzucaro.com/html/

gen:
	hugo

server:
	hugo server -w

combined-css:
	cat static/static/css/normalize.css static/static/css/skeleton.css static/static/css/blog.css > static/static/css/combined.css

css: combined-css
	yuicompressor --type css -o static/static/css/blog.min.css static/static/css/combined.css
	rm static/static/css/combined.css
