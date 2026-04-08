DISCLAIMER ==> I commit .ENV because this project was meant to be only for presentation purpose(I'm not going to do this in prod. I swear), also I plan it to be like the reviewer need to mostly only run build_docker.sh script just to setup&review, but its not done yet you still have to excute the sql yourself to set up the relations, I'll make a pre-seeding function to it later if I'm able to,

WARNING: by runing build_docker.sh will remove your postgresql volume(from my understanding) I do not know if the reviewer need it not to remove or not, but I contemporary need it to prevent conflict with previous build of other project.

my personal note: query params for testing /products, {{base_url}}/products?page=1&limit=5, {{base_url}}/products?q=Go, {{base_url}}/products?category_id=1, {{base_url}}/products?q=Design&category_id=3&sort=price_desc, {{base_url}}/products?sort=newest, {{base_url}}/products?sort=price_asc
