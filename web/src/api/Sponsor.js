import { Request } from "./Base";

const SearchVendorProducts = async (keywords) => {
  const params = new URLSearchParams();
  params.set("q", keywords);

  const query = params.toString();

  return await Request("GET", `/sponsor/vendor/search?${query}`);
};

const GetVendorProduct = async (productID) =>
  await Request("GET", `/sponsor/vendor/products/${productID}`);

const AddVendorProductToCatalog = async (productID) =>
  await Request("POST", `/sponsor/vendor/products/${productID}/add`);

const GetCatalog = async () => await Request("GET", `/sponsor/catalog`);

const GetCatalogProduct = async (productID) =>
  await Request("GET", `/sponsor/catalog/products/${productID}`);

const RemoveCatalogProduct = async (productID) =>
  await Request("GET", `/sponsor/catalog/products/${productID}/remove`);

export {
  SearchVendorProducts,
  GetVendorProduct,
  AddVendorProductToCatalog,
  GetCatalog,
  GetCatalogProduct,
  RemoveCatalogProduct,
};
