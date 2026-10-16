import { Link, useLocation, useNavigate, useSearchParams } from "react-router-dom";
import { useAuth } from "../auth/useAuth";
import { cartCount, useCart } from "../cart/useCart";
import UserMenu from "./UserMenu";

export default function Navbar() {
  const token = useAuth((s) => s.token);
  const items = useCart((s) => s.items);
  const count = cartCount(items);
  const navigate = useNavigate();
  const location = useLocation();
  const [params] = useSearchParams();

  const onSearch = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const q = new FormData(e.currentTarget).get("q")?.toString().trim() ?? "";
    navigate(q ? `/?q=${encodeURIComponent(q)}` : "/");
  };

  return (
    <header className="navbar">
      <div className="nav-inner">
        <Link to="/" className="brand">
          RodriShop
        </Link>

        {location.pathname === "/" && (
          <form className="nav-search" onSubmit={onSearch} role="search">
            <input
              name="q"
              key={params.get("q") ?? ""}
              defaultValue={params.get("q") ?? ""}
              placeholder="Buscar productos"
              aria-label="Buscar productos"
            />
          </form>
        )}

        <nav className="nav-actions">
          <Link to="/carrito" className="cart-link" aria-label="Carrito">
            Carrito
            {count > 0 && <span className="badge">{count}</span>}
          </Link>

          {token ? (
            <UserMenu />
          ) : (
            <>
              <Link to="/login" className="nav-link">
                Iniciar sesión
              </Link>
              <Link to="/register" className="btn btn-small">
                Crear cuenta
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
