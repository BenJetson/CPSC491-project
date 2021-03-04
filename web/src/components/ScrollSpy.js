import { useEffect } from "react";
import { useLocation } from "react-router-dom";

// Implements React Router scroll restoration.
// https://reactrouter.com/web/guides/scroll-restoration

let ScrollSpy = () => {
  const pathName = useLocation();

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathName]);

  return null;
};

export default ScrollSpy;
