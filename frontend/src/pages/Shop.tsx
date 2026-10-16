import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { fetchProducts, type Product } from "../api/shop";
import ProductCard from "../components/ProductCard";

export default function Shop() {
  const [params] = useSearchParams();
  const query = (params.get("q") ?? "").trim().toLowerCase();
  const [products, setProducts] = useState<Product[] | null>(null);
  const [error, setError] = useState(false);

  useEffect(() => {
    fetchProducts()
      .then(setProducts)
      .catch(() => setError(true));
  }, []);

  const visible = useMemo(() => {
    if (!products) return [];
    if (!query) return products;
    return products.filter(
      (p) => p.name.toLowerCase().includes(query) || p.description?.toLowerCase().includes(query),
    );
  }, [products, query]);

  return (
    <>
      <section className="hero">
        <h1>Bienvenido a RodriShop</h1>
        <p>Productos seleccionados, envío a todo el país y pago seguro.</p>
      </section>

      <div className="section-head">
        <h2>{query ? `Resultados para “${params.get("q")}”` : "Productos"}</h2>
        {products && <span className="muted">{visible.length} productos</span>}
      </div>

      {error && <p className="msg msg-error">No pudimos cargar el catálogo. Probá de nuevo en un rato.</p>}
      {!error && products === null && <p className="muted">Cargando productos...</p>}
      {products && visible.length === 0 && (
        <p className="muted">
          {query ? "No encontramos productos con esa búsqueda." : "Todavía no hay productos disponibles."}
        </p>
      )}

      <div className="product-grid">
        {visible.map((p) => (
          <ProductCard key={p.id} product={p} />
        ))}
      </div>
    </>
  );
}
