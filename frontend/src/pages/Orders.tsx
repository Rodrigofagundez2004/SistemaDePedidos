import { useCallback, useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import {
  fetchOrderDetail,
  fetchOrders,
  formatMoney,
  payOrder,
  STATUS_LABELS,
  type Order,
  type OrderItem,
} from "../api/shop";

interface Notice {
  type: "ok" | "error";
  text: string;
}

export default function Orders() {
  const location = useLocation();
  const [orders, setOrders] = useState<Order[] | null>(null);
  const [details, setDetails] = useState<Record<string, OrderItem[]>>({});
  const [payingId, setPayingId] = useState<string | null>(null);
  const [notice, setNotice] = useState<Notice | null>(
    (location.state as { notice?: Notice } | null)?.notice ?? null,
  );
  const [error, setError] = useState(false);

  const load = useCallback(() => {
    fetchOrders()
      .then((list) => {
        setOrders(list);
        setError(false);
        list.forEach((o) => {
          fetchOrderDetail(o.id)
            .then((items) => setDetails((d) => ({ ...d, [o.id]: items })))
            .catch(() => undefined);
        });
      })
      .catch(() => setError(true));
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const retry = async (orderId: string) => {
    setPayingId(orderId);
    setNotice(null);
    try {
      const payment = await payOrder(orderId);
      setNotice(
        payment.status === "succeeded"
          ? { type: "ok", text: "¡Pago aprobado! Tu pedido quedó pagado." }
          : { type: "error", text: "El pago volvió a ser rechazado. Podés intentarlo más tarde." },
      );
    } catch (err: any) {
      setNotice({ type: "error", text: err.response?.data?.error || "No se pudo procesar el pago." });
    } finally {
      setPayingId(null);
      load();
    }
  };

  return (
    <>
      <div className="section-head">
        <h1>Mis pedidos</h1>
      </div>

      {notice && <p className={`msg ${notice.type === "ok" ? "msg-ok" : "msg-error"}`}>{notice.text}</p>}
      {error && <p className="msg msg-error">No pudimos cargar tus pedidos.</p>}
      {!error && orders === null && <p className="muted">Cargando pedidos...</p>}
      {orders && orders.length === 0 && <p className="muted">Todavía no hiciste ningún pedido.</p>}

      <ul className="order-list">
        {orders?.map((o) => (
          <li className="order-card" key={o.id}>
            <div className="order-head">
              <div>
                <strong>Pedido #{o.id.slice(0, 8).toUpperCase()}</strong>
                <span className="muted order-date">{new Date(o.created_at).toLocaleDateString("es-UY")}</span>
              </div>
              <span className={`status status-${o.status}`}>{STATUS_LABELS[o.status] ?? o.status}</span>
            </div>
            {details[o.id] && (
              <ul className="order-items">
                {details[o.id].map((it) => (
                  <li key={it.id}>
                    {it.quantity} × {it.product_name}
                    <span className="muted"> — {formatMoney(it.subtotal_cents, o.currency)}</span>
                  </li>
                ))}
              </ul>
            )}
            <div className="order-foot">
              <strong>Total: {formatMoney(o.total_cents, o.currency)}</strong>
              {(o.status === "pending" || o.status === "failed") && (
                <button type="button" className="btn btn-small" onClick={() => retry(o.id)} disabled={payingId !== null}>
                  {payingId === o.id ? "Procesando..." : o.status === "failed" ? "Reintentar pago" : "Pagar ahora"}
                </button>
              )}
            </div>
          </li>
        ))}
      </ul>
    </>
  );
}
