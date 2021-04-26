const ZeroPad = (n, width) => {
  // Inspired by: https://stackoverflow.com/a/10073788/6127099

  const z = "0";

  n = n + "";
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
};

const FormatMoney = (money) => {
  const dollars = Math.floor(money / 100);
  const cents = money % 100;

  return `$${dollars}.${ZeroPad(cents, 2)}`;
};

export { ZeroPad, FormatMoney };
