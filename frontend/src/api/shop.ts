import api from "./client";

// Contratos según el backend actual:
//  - GET /orders          -> { data: Order[], total, page, limit }
//  - GET /orders/:id      -> { order: Order, items: OrderItem[] }
//  - POST /orders         -> Order  (body: { items, idempotency_key })
//  - POST /payments       -> Payment (body: { order_id, idempotency_key })
//  - GET /products        -> { data: Product[] }  (todavía no existe en el backend actual)

export interface Product {
  id: string;
  name: string;
  description: string;
  price_cents: number;
  currency: string;
  stock: number;
  image_url: string;
  is_active: boolean;
}

export interface Order {
  id: string;
  status: string;
  total_cents: number;
  currency: string;
  created_at: string;
}

export interface OrderItem {
  id: string;
  product_id: string;
  product_name: string;
  unit_price_cents: number;
  quantity: number;
  subtotal_cents: number;
}

export interface Payment {
  id: string;
  order_id: string;
  amount_cents: number;
  currency: string;
  status: string;
}

export const STATUS_LABELS: Record<string, string> = {
  pending: "Pendiente de pago",
  paid: "Pagado",
  failed: "Pago rechazado",
};

export function formatMoney(cents: number, currency: string): string {
  return new Intl.NumberFormat("es-UY", { style: "currency", currency }).format(cents / 100);
}

export async function fetchProducts(): Promise<Product[]> {
  const { data } = await api.get<{ data: Product[] }>("/products");
  return data.data;
}

export async function createOrder(items: { product_id: string; quantity: number }[]): Promise<Order> {
  const { data } = await api.post<Order>("/orders", {
    items,
    idempotency_key: crypto.randomUUID(),
  });
  return data;
}

export async function fetchOrders(): Promise<Order[]> {
  const { data } = await api.get<{ data: Order[] }>("/orders", { params: { limit: 50 } });
  return data.data;
}

export async function fetchOrderDetail(id: string): Promise<OrderItem[]> {
  const { data } = await api.get<{ items: OrderItem[] }>(`/orders/${id}`);
  return data.items ?? [];
}

export async function payOrder(orderId: string): Promise<Payment> {
  // Clave nueva por intento: el backend la usa para no cobrar dos veces.
  const { data } = await api.post<Payment>("/payments", {
    order_id: orderId,
    idempotency_key: crypto.randomUUID(),
  });
  return data;
}
