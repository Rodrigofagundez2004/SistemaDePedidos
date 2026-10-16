import { useEffect, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../auth/useAuth";

export default function UserMenu() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const onClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", onClick);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onClick);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  if (!user) return null;

  const initial = (user.full_name || user.email).trim().charAt(0).toUpperCase();

  const handleLogout = () => {
    setOpen(false);
    logout();
    navigate("/");
  };

  return (
    <div className="user-menu" ref={ref}>
      <button
        type="button"
        className="avatar"
        aria-label="Menú de usuario"
        aria-expanded={open}
        onClick={() => setOpen((o) => !o)}
      >
        {initial}
      </button>

      {open && (
        <div className="menu-dropdown" role="menu">
          <div className="menu-head">
            <strong>{user.full_name}</strong>
            <span className="muted">{user.email}</span>
          </div>
          <Link to="/cuenta" role="menuitem" onClick={() => setOpen(false)}>
            Mi cuenta
          </Link>
          <Link to="/pedidos" role="menuitem" onClick={() => setOpen(false)}>
            Mis pedidos
          </Link>
          <div className="menu-divider" />
          <button type="button" role="menuitem" className="menu-logout" onClick={handleLogout}>
            Cerrar sesión
          </button>
        </div>
      )}
    </div>
  );
}
