import { useCallback, useEffect, useState } from "react";
import {
  fetchOrderDetail,
  fetchOrders,
  formatMoney,
  payOrder,
  STATUS_LABELS,
  type Order,
  type OrderItem,
} from "../api/shop";

interface Props {
  refreshKey: number;
}

export default function OrdersSection({ refreshKey }: Props) {
  const [orders, setOrders] = useState<Order[] | null>(null);
  const [details, setDetails] = useState<Record<string, OrderItem[]>>({});
  const [payingId, setPayingId] = useState<string | null>(null);
  const [message, setMessage] = useState<{ type: "ok" | "error"; text: string } | null>(null);
  const [loadError, setLoadError] = useState(false);

  const load = useCallback(() => {
    fetchOrders()
      .then((list) => {
        setOrders(list);
        setLoadError(false);
      })
      .catch(() => setLoadError(true));
  }, []);

  useEffect(() => {
    load();
  }, [load, refreshKey]);

  // Carga los ítems de cada pedido la primera vez (la lista no los trae).
  useEffect(() => {
    if (!orders) return;
    orders.forEach((o) => {
      if (!details[o.id]) {
        fetchOrderDetail(o.id)
          .then((items) => setDetails((d) => ({ ...d, [o.id]: items })))
          .catch(() => undefined);
      }
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [orders]);

  const handlePay = async (orderId: string) => {
    setPayingId(orderId);
    setMessage(null);
    try {
      const payment = await payOrder(orderId);
      setMessage(
        payment.status === "succeeded"
          ? { type: "ok", text: "¡Pago aprobado! Tu pedido quedó pagado." }
          : { type: "error", text: "El pago fue rechazado. Podés intentarlo de nuevo." },
      );
    } catch (err: any) {
      setMessage({ type: "error", text: err.response?.data?.error || "No se pudo procesar el pago." });
    } finally {
      setPayingId(null);
      load();
    }
  };

  if (loadError) return <p className="msg msg-error">No pudimos cargar tus pedidos.</p>;
  if (orders === null) return <p className="muted">Cargando pedidos...</p>;
  if (orders.length === 0) return <p className="muted">Todavía no tenés pedidos.</p>;

  return (
    <>
      <ul className="order-list">
        {orders.map((o) => (
          <li className="order-card" key={o.id}>
            <div className="order-head">
              <strong>Pedido #{o.id.slice(0, 8)}</strong>
              <span className={`status status-${o.status}`}>{STATUS_LABELS[o.status] ?? o.status}</span>
            </div>
            {details[o.id] && (
              <ul className="order-items">
                {details[o.id].map((it) => (
                  <li key={it.id}>
                    {it.quantity} × {it.product_name} — {formatMoney(it.subtotal_cents, o.currency)}
                  </li>
                ))}
              </ul>
            )}
            <div className="order-foot">
              <span>Total: {formatMoney(o.total_cents, o.currency)}</span>
              {(o.status === "pending" || o.status === "failed") && (
                <button type="button" className="btn btn-small" onClick={() => handlePay(o.id)} disabled={payingId !== null}>
                  {payingId === o.id ? "Procesando..." : o.status === "failed" ? "Reintentar pago" : "Pagar"}
                </button>
              )}
            </div>
          </li>
        ))}
      </ul>
      {message && <p className={`msg ${message.type === "ok" ? "msg-ok" : "msg-error"}`}>{message.text}</p>}
    </>
  );
}
