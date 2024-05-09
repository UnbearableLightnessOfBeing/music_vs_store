INSERT INTO categories (name, slug, img_url) VALUES ('Клавишные инструменты', 'keyboards', '/assets/images/categories/keyboards.jpg');
INSERT INTO categories (name, slug, img_url) VALUES ('Гитары', 'guitars', '/assets/images/categories/guitars.jpg');
INSERT INTO categories (name, slug, img_url) VALUES ('Духовые инструменты', 'wind-instruments', '/assets/images/categories/wind-instruments.jpg');
INSERT INTO categories (name, slug, img_url) VALUES ('Ударные инструменты', 'percussion', '/assets/images/categories/percussion.jpg');
INSERT INTO categories (name, slug, img_url) VALUES ('Звуковое оборудование', 'sound-equipment', '/assets/images/categories/sound-equipment.jpg');
INSERT INTO categories (name, slug, img_url) VALUES ('Микрофоны', 'mics', '/assets/images/categories/mics.jpg');
INSERT INTO categories (name, slug, img_url) VALUES ('Разъемы и кабели', 'cables', '/assets/images/categories/cables.jpg');

insert into labels (name) values ('Honor');
insert into labels (name) values ('Huawey');
insert into labels (name) values ('LG');

insert into products (name, price_int, label_id, images, description, in_stock) values ('Alhambra', '1200', '1','{ "/assets/images/products/alhambra.jpg", "/assets/images/products/alhambra_1.jpg", "/assets/images/products/alhambra_2.jpg" }'::varchar[], 'cool guitar', true);
insert into products (name, price_int, label_id, images, description, in_stock) values ('ALMANSA', '900', '2','{ "/assets/images/products/almansa.jpg", "/assets/images/products/almansa_1.jpg", "/assets/images/products/almansa_2.jpg", "/assets/images/products/almansa_3.jpg", "/assets/images/products/almansa_4.jpg", "/assets/images/products/almansa_5.jpg", "/assets/images/products/almansa_6.jpg" }'::varchar[], 'ALMANSA 402 Cedro – 6-струнная полноразмерная классическая гитара.
Модель из серии гитар (Estudio). Верхняя дека из массива кедра, корпус из слоеного красного дерева, гриф из красного дерева с накладкой из индийского палисандра. Мензура 650 мм, ширина верхнего порожка 52 мм. Глянцевая отделка.
Произведена в Испании.', true);
insert into product_categories (product_id, category_id) values ('1', '2');
insert into product_categories (product_id, category_id) values ('2', '2');

insert into products (name, price_int, label_id, description, in_stock) values ('KAWAI CR-40 TRANSPARENCY', '4000', '3', 'cool piano', true);
insert into product_categories (product_id, category_id) values ('3', '1');

insert into countries (name) values ('Россия');
insert into countries (name) values ('Беларусь');
insert into countries (name) values ('Украина');
insert into countries (name) values ('Казахстан');

insert into delivery_methods (name) values ('ТК Деловые линии');
insert into delivery_methods (name) values ('ТК СДЕК');
insert into delivery_methods (name) values ('Курьером');
insert into delivery_methods (name) values ('Самовывоз');

insert into payment_methods (name) values ('Банковской картой');
insert into payment_methods (name) values ('Наличные');

-- pass: admin
INSERT INTO users (username, email, is_admin, password) VALUES ('admin', 'admin@admin.ru', 'true', '$2a$10$DS/dEObNUtdY4Q6LCdbSf.FjRE3y87tB0pC9bwSaiVADQK5tGoHEm');
-- pass: user
INSERT INTO users (username, email, is_admin, password) VALUES ('user', 'user@user.com', 'false', '$2a$10$Rx6oImAgXCqrlKlz1nRNlOnLLPtVDksvevkmMazB0XCMGkcJWHHTi');
