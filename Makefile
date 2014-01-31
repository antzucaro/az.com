deploy:
	rsync --size-only -av public/ azucaro@antzucaro.com:/home/azucaro/antzucaro.com/public/

gen:
	hugo

server:
	hugo server -w
