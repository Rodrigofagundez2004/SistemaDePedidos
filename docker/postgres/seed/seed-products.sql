-- Productos de demo (20). Se puede correr varias veces sin duplicar.
-- Se ejecuta contra products_db, después de que el backend creó las tablas.
INSERT INTO products (name, description, price_cents, currency, stock, image_url, is_active, created_at, updated_at)
SELECT v.name, v.description, v.price_cents, 'USD', v.stock, v.image_url, true, now(), now()
FROM (VALUES
  ('Remera negra',            'Algodón 100%, corte recto',              2500, 20, 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=600'),
  ('Remera blanca',           'Algodón peinado, cuello redondo',        2200, 15, 'https://images.unsplash.com/photo-1583743814966-8936f5b7be1a?w=600'),
  ('Remera estampada',        'Algodón, estampa frontal',               2800, 12, ''),
  ('Camisa de lino',          'Lino liviano, ideal para verano',        4500,  8, ''),
  ('Buzo gris',               'Frisa suave con capucha',                4000,  5, 'https://images.unsplash.com/photo-1556821840-3a63f95609a7?w=600'),
  ('Buzo negro con capucha',  'Frisa gruesa, bolsillo canguro',         4800,  7, ''),
  ('Campera de jean',         'Clásica, corte regular',                 5500,  8, 'https://images.unsplash.com/photo-1576995853123-5a10305d93c0?w=600'),
  ('Campera rompeviento',     'Liviana e impermeable al viento',        6200,  6, ''),
  ('Gorra',                   'Ajustable, unisex',                      1000, 10, 'https://images.unsplash.com/photo-1588850561407-ed78c334e67a?w=600'),
  ('Gorra trucker',           'Malla trasera y visera curva',           1500,  9, ''),
  ('Medias pack x3',          'Algodón, uso diario',                     800, 30, ''),
  ('Pantalón jogger',         'Algodón con puño elástico',              3500, 14, ''),
  ('Jean clásico',            'Denim, corte recto',                     4900, 11, ''),
  ('Short de baño',           'Secado rápido, con forro',               2600, 13, ''),
  ('Bermuda de lino',         'Cintura elástica, livianas',             3200, 10, ''),
  ('Vestido floral',          'Viscosa, largo midi',                    4200,  6, ''),
  ('Pollera midi',            'Plisada, cintura alta',                  3600,  9, ''),
  ('Zapatillas urbanas',      'Suela de goma, uso diario',              8900,  4, ''),
  ('Mochila urbana',          'Compartimento para notebook',            3900, 12, ''),
  ('Cinturón de cuero',       'Hebilla metálica, 3 cm de ancho',        1800, 18, '')
) AS v(name, description, price_cents, stock, image_url)
WHERE NOT EXISTS (SELECT 1 FROM products p WHERE p.name = v.name);