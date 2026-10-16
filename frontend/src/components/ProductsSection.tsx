import { useEffect, useState } from "react";
import { createOrder, fetchProducts, formatMoney, type Product } from "../api/shop";

interface Props {
  onOrderCreated: () => void;
}

export default function ProductsSection({ onOrderCreated }: Props) {
  const [products, setProducts] = useState<Product[] | null>(null);
  const [unavailable, setUnavailable] = useState(false);
  const [cart, setCart] = useState<Record<string, number>>({});
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState<{ type: "ok" | "error"; text: string } | null>(null);

  useEffect(() => {
    fetchProducts()
      .then(setProducts)
      .catch(() => setUnavailable(true));
  }, []);

  const setQty = (p: Product, qty: number) => {
    const clamped = Math.max(0, Math.min(qty, p.stock));
    setCart((c) => ({ ...c, [p.id]: clamped }));
  };

  const cartLines = (products ?? [])
    .filter((p) => (cart[p.id] ?? 0) > 0)
    .map((p) => ({ product: p, quantity: cart[p.id] }));
  const total = cartLines.reduce((sum, l) => sum + l.product.price_cents * l.quantity, 0);
  const currency = products?.[0]?.currency ?? "USD";

  const handleOrder = async () => {
    setSubmitting(true);
    setMessage(null);
    try {
      await createOrder(cartLines.map((l) => ({ product_id: l.product.id, quantity: l.quantity })));
      setCart({});
      setMessage({ type: "ok", text: "Pedido creado. Podés pagarlo desde Mis pedidos." });
      setProducts(await fetchProducts());
      onOrderCreated();
    } catch (err: any) {
      setMessage({ type: "error", text: err.response?.data?.error || "No se pudo crear el pedido." });
    } finally {
      setSubmitting(false);
    }
  };

  if (unavailable) {
    return <p className="muted">El catálogo todavía no está disponible.</p>;
  }
  if (products === null) {
    return <p className="muted">Cargando productos...</p>;
  }
  if (products.length === 0) {
    return <p className="muted">Todavía no hay productos cargados.</p>;
  }

  return (
    <>
      <div className="product-grid">
        {products.map((p) => {
          const qty = cart[p.id] ?? 0;
          const outOfStock = p.stock === 0;
          return (
            <div className="product-card" key={p.id}>
              <h4>{p.name}</h4>
              {p.description && <p className="muted">{p.description}</p>}
              <div className="product-price">{formatMoney(p.price_cents, p.currency)}</div>
              <div className="product-stock">{outOfStock ? "Sin stock" : `${p.stock} disponibles`}</div>
              <div className="qty-row">
                <button type="button" className="btn btn-small" onClick={() => setQty(p, qty - 1)} disabled={qty === 0}>
                  −
                </button>
                <span>{qty}</span>
                <button type="button" className="btn btn-small" onClick={() => setQty(p, qty + 1)} disabled={outOfStock || qty >= p.stock}>
                  +
                </button>
              </div>
            </div>
          );
        })}
      </div>

      <div className="cart-bar">
        <span>{cartLines.length === 0 ? "Elegí productos para armar tu pedido" : `Total: ${formatMoney(total, currency)}`}</span>
        <button type="button" className="btn" onClick={handleOrder} disabled={cartLines.length === 0 || submitting}>
          {submitting ? "Creando pedido..." : "Crear pedido"}
        </button>
      </div>

      {message && <p className={`msg ${message.type === "ok" ? "msg-ok" : "msg-error"}`}>{message.text}</p>}
    </>
  );
}
