import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { createOrder, formatMoney, payOrder } from "../api/shop";
import { useAuth } from "../auth/useAuth";
import { cartTotal, useCart } from "../cart/useCart";

export default function Cart() {
  const items = useCart((s) => s.items);
  const setQuantity = useCart((s) => s.setQuantity);
  const remove = useCart((s) => s.remove);
  const clear = useCart((s) => s.clear);
  const token = useAuth((s) => s.token);
  const navigate = useNavigate();
  const [busy, setBusy] = useState(false);
  const [message, setMessage] = useState<{ type: "ok" | "error"; text: string } | null>(null);

  const currency = items[0]?.currency ?? "USD";
  const total = cartTotal(items);

  const checkout = async () => {
    if (!token) {
      navigate("/login", { state: { from: "/carrito" } });
      return;
    }
    setBusy(true);
    setMessage(null);
    try {
      const order = await createOrder(items.map((i) => ({ product_id: i.product_id, quantity: i.quantity })));
      const payment = await payOrder(order.id);
      if (payment.status === "succeeded") {
        clear();
        navigate("/pedidos", { state: { notice: { type: "ok", text: "¡Gracias por tu compra! Tu pago fue aprobado." } } });
      } else {
        navigate("/pedidos", {
          state: { notice: { type: "error", text: "El pago fue rechazado. Tu pedido quedó guardado: podés reintentar el pago." } },
        });
      }
    } catch (err: any) {
      setMessage({ type: "error", text: err.response?.data?.error || "No pudimos completar la compra. Revisá el stock y probá de nuevo." });
    } finally {
      setBusy(false);
    }
  };

  if (items.length === 0) {
    return (
      <div className="empty-state">
        <h1>Tu carrito está vacío</h1>
        <p className="muted">Explorá el catálogo y agregá los productos que te gusten.</p>
        <Link to="/" className="btn">
          Ver productos
        </Link>
      </div>
    );
  }

  return (
    <div className="cart-layout">
      <section className="cart-items">
        <h1>Carrito</h1>
        {items.map((i) => (
          <div className="cart-line" key={i.product_id}>
            <div className="cart-thumb">
              {i.image_url ? <img src={i.image_url} alt={i.name} /> : <span>{i.name.charAt(0)}</span>}
            </div>
            <div className="cart-info">
              <strong>{i.name}</strong>
              <span className="muted">{formatMoney(i.price_cents, i.currency)} c/u</span>
              <button type="button" className="link-btn" onClick={() => remove(i.product_id)}>
                Quitar
              </button>
            </div>
            <div className="qty-row">
              <button type="button" className="btn btn-small" onClick={() => setQuantity(i.product_id, i.quantity - 1)}>
                −
              </button>
              <span>{i.quantity}</span>
              <button
                type="button"
                className="btn btn-small"
                disabled={i.quantity >= i.stock}
                onClick={() => setQuantity(i.product_id, i.quantity + 1)}
              >
                +
              </button>
            </div>
            <strong className="line-total">{formatMoney(i.price_cents * i.quantity, i.currency)}</strong>
          </div>
        ))}
      </section>

      <aside className="cart-summary">
        <h2>Resumen</h2>
        <div className="summary-row">
          <span>Productos</span>
          <span>{items.reduce((n, i) => n + i.quantity, 0)}</span>
        </div>
        <div className="summary-row total">
          <span>Total</span>
          <span>{formatMoney(total, currency)}</span>
        </div>
        <button type="button" className="btn btn-block" onClick={checkout} disabled={busy}>
          {busy ? "Procesando compra..." : token ? "Finalizar compra" : "Iniciá sesión para comprar"}
        </button>
        {message && <p className="msg msg-error">{message.text}</p>}
      </aside>
    </div>
  );
}
