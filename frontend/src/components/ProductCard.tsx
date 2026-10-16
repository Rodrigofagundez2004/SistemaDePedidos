import { useState } from "react";
import { formatMoney, type Product } from "../api/shop";
import { useCart } from "../cart/useCart";

export default function ProductCard({ product }: { product: Product }) {
  const add = useCart((s) => s.add);
  const inCart = useCart((s) => s.items.find((i) => i.product_id === product.id)?.quantity ?? 0);
  const [imgFailed, setImgFailed] = useState(false);
  const outOfStock = product.stock <= 0;
  const reachedLimit = inCart >= product.stock;

  return (
    <article className="product-card">
      <div className="product-img">
        {product.image_url && !imgFailed ? (
          <img src={product.image_url} alt={product.name} onError={() => setImgFailed(true)} />
        ) : (
          <span>{product.name.charAt(0).toUpperCase()}</span>
        )}
      </div>
      <div className="product-body">
        <h3>{product.name}</h3>
        {product.description && <p className="product-desc">{product.description}</p>}
        <div className="product-foot">
          <span className="product-price">{formatMoney(product.price_cents, product.currency)}</span>
          {outOfStock ? (
            <span className="product-stock out">Sin stock</span>
          ) : (
            <span className="product-stock">Stock: {product.stock}</span>
          )}
        </div>
        <button
          type="button"
          className="btn btn-block"
          disabled={outOfStock || reachedLimit}
          onClick={() =>
            add({
              product_id: product.id,
              name: product.name,
              price_cents: product.price_cents,
              currency: product.currency,
              image_url: product.image_url,
              stock: product.stock,
            })
          }
        >
          {outOfStock ? "No disponible" : reachedLimit ? "Máximo en el carrito" : inCart ? `Agregar otro (${inCart} en carrito)` : "Agregar al carrito"}
        </button>
      </div>
    </article>
  );
}
