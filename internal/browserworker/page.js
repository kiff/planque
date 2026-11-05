// The page as facts, for planque: fields with their kind, controls, links,
// a payment frame, a checkout's total. Secret values never leave here.
// Each element gets a stable ref (data-planque-ref) later steps name.
(() => {
  const MAX_TEXT = 65536;
  const n = window.__planqueRefs || (window.__planqueRefs = {f: 0, c: 0, l: 0});
  const ref = (el, k) => {
    if (!el.dataset.planqueRef) el.dataset.planqueRef = k + (++n[k]);
    return el.dataset.planqueRef;
  };
  const visible = el => {
    const s = getComputedStyle(el);
    if (s.visibility === 'hidden' || s.display === 'none') return false;
    const r = el.getBoundingClientRect();
    return r.width > 0 || r.height > 0;
  };
  const clean = s => (s || '').replace(/\s+/g, ' ').trim().slice(0, 120);
  const name = el => clean(
    el.getAttribute('aria-label') ||
    (el.labels && el.labels[0] && el.labels[0].textContent) ||
    (el.getAttribute('aria-labelledby') && (document.getElementById(el.getAttribute('aria-labelledby')) || {}).textContent) ||
    el.getAttribute('placeholder') || el.getAttribute('title') ||
    (el.querySelector && el.querySelector('img[alt]') && el.querySelector('img[alt]').alt) ||
    (el.tagName === 'INPUT' ? el.value && ['submit', 'button', 'image'].includes(el.type) && el.value : el.textContent) ||
    el.getAttribute('name') || '');
  const kindOf = el => {
    const ac = (el.getAttribute('autocomplete') || '').toLowerCase();
    const t = (el.type || 'text').toLowerCase();
    const nm = ((el.name || '') + ' ' + (el.id || '') + ' ' + (el.getAttribute('data-elements') || '')).toLowerCase();
    if (t === 'password') return 'password';
    if (/cc-number|cc-csc|cc-exp/.test(ac) || /card.?number|cardnumber|\bcvc\b|\bcvv\b|\bcsc\b|security.?code|cc[-_]?num|pan\b/.test(nm)) return 'card';
    if (ac === 'one-time-code' || /\botp\b|one.?time|verification.?code|\btotp\b|\b2fa\b|\bmfa\b|sms.?code|auth.?code/.test(nm)) return 'code';
    if (el.tagName === 'SELECT') return 'select';
    if (el.tagName === 'TEXTAREA') return 'textarea';
    if (['checkbox', 'radio', 'email', 'tel', 'number', 'search', 'url', 'date'].includes(t)) return t;
    return 'text';
  };
  const secret = k => k === 'password' || k === 'card' || k === 'code';
  const skip = t => ['hidden', 'submit', 'button', 'image', 'reset', 'file'].includes(t);
  const fields = [...document.querySelectorAll('input, select, textarea')]
    .filter(el => !skip((el.type || '').toLowerCase()) && !el.disabled && visible(el))
    .map(el => {
      const k = kindOf(el);
      const f = {ref: ref(el, 'f'), label: name(el), kind: k};
      if (secret(k)) {
        f.masked = true;
      } else if (el.tagName === 'SELECT') {
        f.value = clean(el.value);
        f.options = [...el.options].map(o => clean(o.textContent)).slice(0, 50);
      } else if (k === 'checkbox' || k === 'radio') {
        f.value = el.checked ? 'on' : '';
      } else {
        f.value = (el.value || '').slice(0, 500);
      }
      return f;
    });
  const formOf = el => el.form || el.closest('form');
  const cardForm = el => {
    const f = formOf(el);
    return !!f && [...f.querySelectorAll('input')].some(i => kindOf(i) === 'card');
  };
  const orderPath = /checkout|order|pay|purchase|buy|kauf|bestell|commande|achat|pedido|compra|acquist|koop|betal/i;
  const orderAttrs = /place.?order|submit.?order|checkout.?submit|buy.?now|one.?click|pay.?now|purchase|place-order|btn-buy|btn-pay|order-button|buy-button/;
  const placesOrder = el => {
    if (cardForm(el)) return true;
    const f = formOf(el);
    const action = el.getAttribute('formaction') || (f && f.getAttribute('action')) || '';
    if (action && orderPath.test(action.split('?')[0])) return true;
    if (el.closest('[itemtype*="BuyAction"], [itemtype*="PayAction"], [itemtype*="OrderAction"]')) return true;
    const attrs = ((el.id || '') + ' ' + (el.className && el.className.baseVal !== undefined ? '' : el.className || '') + ' ' +
      (el.getAttribute('name') || '') + ' ' + (el.getAttribute('data-testid') || '') + ' ' + (el.getAttribute('data-action') || '')).toLowerCase();
    if (orderAttrs.test(attrs)) return true;
    if (el.closest('[data-stripe], [id*="paypal"], [class*="paypal"], [id*="stripe"], [class*="stripe-button"], [class*="apple-pay"], [class*="google-pay"]')) return true;
    return false;
  };
  const controls = [...document.querySelectorAll('button, input[type=submit], input[type=button], input[type=image], [role=button], summary')]
    .filter(el => !el.disabled && visible(el))
    .map(el => {
      const c = {ref: ref(el, 'c'), label: name(el)};
      const t = (el.type || '').toLowerCase();
      if (formOf(el) && ((el.tagName === 'BUTTON' && t !== 'button') || t === 'submit' || t === 'image')) c.submits = true;
      if (placesOrder(el)) c.places_order = true;
      return c;
    });
  const links = [...document.querySelectorAll('a[href]')]
    .filter(el => visible(el) && /^https?:/.test(el.href))
    .slice(0, 200)
    .map(el => ({ref: ref(el, 'l'), label: name(el), url: el.href}));
  const paymentFrame = [...document.querySelectorAll('iframe')]
    .some(f => /stripe|paypal|adyen|braintree|checkout\.com|klarna|mollie|squareup|worldpay|payment|card/i.test((f.src || '') + ' ' + (f.name || '') + ' ' + (f.title || '')));

  // The checkout's total: one amount next to a "total" word, in one
  // currency. Two different candidates, or none, is no total read.
  const codes = 'EUR|USD|GBP|CHF|SEK|NOK|DKK|PLN|CZK|HUF|JPY|CAD|AUD|BRL|MXN';
  const money = new RegExp('(?:(' + codes + ')|([€$£]))\\s?(\\d{1,3}(?:[.,\\s]\\d{3})*(?:[.,]\\d{2})?)|(\\d{1,3}(?:[.,\\s]\\d{3})*(?:[.,]\\d{2})?)\\s?(?:(' + codes + ')|([€$£]))', 'g');
  const symbol = {'€': 'EUR', '$': 'USD', '£': 'GBP'};
  const normalize = s => {
    s = s.replace(/\s/g, '');
    const m = s.match(/^(.*?)([.,])(\d{2})$/);
    let whole = s, cents = '00';
    if (m) { whole = m[1]; cents = m[3]; }
    whole = whole.replace(/[.,]/g, '');
    return whole + '.' + cents;
  };
  const amounts = text => {
    const out = [];
    for (const m of text.matchAll(money)) {
      const cur = m[1] || symbol[m[2]] || m[5] || symbol[m[6]];
      const amt = m[3] || m[4];
      if (cur && amt) out.push({amount: normalize(amt), currency: cur});
    }
    return out;
  };
  const totalWord = /\b(total|totale|gesamt(?:betrag|summe|preis)?|summe|totaal|montant total|importe total|a pagar|à payer|zu zahlen|da pagare|te betalen|to pay|amount due)\b/i;
  const strong = /grand total|order total|total to pay|to pay|a pagar|à payer|zu zahlen|da pagare|te betalen|importe total|montant total|gesamtbetrag|amount due/i;
  const own = el => [...el.childNodes].filter(c => c.nodeType === 3).map(c => c.textContent).join(' ');
  const found = [];
  for (const el of document.body ? document.body.querySelectorAll('*') : []) {
    if (['SCRIPT', 'STYLE', 'INPUT', 'TEXTAREA', 'SELECT'].includes(el.tagName) || !visible(el)) continue;
    const text = clean(el.textContent);
    if (text.length > 160 || !totalWord.test(text)) continue;
    let candidates = amounts(text);
    if (candidates.length === 0 && el.parentElement) candidates = amounts(clean(el.parentElement.textContent).slice(0, 240));
    if (candidates.length === 0 && el.nextElementSibling) candidates = amounts(clean(el.nextElementSibling.textContent));
    if (candidates.length === 0) continue;
    // A row that names several amounts (a subtotal and a total) gives the last.
    const a = candidates[candidates.length - 1];
    found.push({...a, strong: strong.test(text)});
  }
  const distinct = list => [...new Set(list.map(a => a.amount + ' ' + a.currency))];
  let total = null;
  const strongOnes = found.filter(a => a.strong);
  if (distinct(strongOnes).length === 1) total = strongOnes[0];
  else if (distinct(found).length === 1) total = found[0];
  const checkout = total ? {lines: [], total: total.amount, currency: total.currency} : null;
  return JSON.stringify({
    url: location.href, title: clean(document.title),
    text: (document.body ? document.body.innerText : '').slice(0, MAX_TEXT),
    fields, controls, links, payment_frame: paymentFrame, checkout,
  });
})()
