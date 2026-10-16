import { create } from "zustand";
import { persist } from "zustand/middleware";

export interface CartItem {
  product_id: string;
  name: string;
  price_cents: number;
  currency: string;
  image_url: string;
  stock: number; // stock conocido al momento de agregar; el backend lo valida al comprar
  quantity: number;
}

interface CartState {
  items: CartItem[];
  add: (item: Omit<CartItem, "quantity">) => void;
  setQuantity: (productId: string, quantity: number) => void;
  remove: (productId: string) => void;
  clear: () => void;
}

export const useCart = create<CartState>()(
  persist(
    (set) => ({
      items: [],
      add: (item) =>
        set((state) => {
          const existing = state.items.find((i) => i.product_id === item.product_id);
          if (existing) {
            const quantity = Math.min(existing.quantity + 1, item.stock);
            return {
              items: state.items.map((i) =>
                i.product_id === item.product_id ? { ...i, stock: item.stock, quantity } : i,
              ),
            };
          }
          if (item.stock <= 0) return state;
          return { items: [...state.items, { ...item, quantity: 1 }] };
        }),
      setQuantity: (productId, quantity) =>
        set((state) => ({
          items: state.items
            .map((i) =>
              i.product_id === productId
                ? { ...i, quantity: Math.max(0, Math.min(quantity, i.stock)) }
                : i,
            )
            .filter((i) => i.quantity > 0),
        })),
      remove: (productId) =>
        set((state) => ({ items: state.items.filter((i) => i.product_id !== productId) })),
      clear: () => set({ items: [] }),
    }),
    { name: "rodrishop-cart" },
  ),
);

export function cartCount(items: CartItem[]): number {
  return items.reduce((n, i) => n + i.quantity, 0);
}

export function cartTotal(items: CartItem[]): number {
  return items.reduce((sum, i) => sum + i.price_cents * i.quantity, 0);
}
