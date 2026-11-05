// One step on the page: fill a field, select an option, click a control
// or a link, by the ref page.js gave it. Returns "ok" or why not.
((step) => {
  const ref = step.fill || step.select || step.click;
  const el = document.querySelector('[data-planque-ref="' + ref + '"]');
  if (!el) return 'ref_unknown';
  const fire = (name) => el.dispatchEvent(new Event(name, {bubbles: true}));
  if (step.fill !== undefined && step.fill !== '') {
    const t = (el.type || '').toLowerCase();
    if (t === 'checkbox' || t === 'radio') {
      el.checked = step.value !== '' && step.value !== 'off';
      fire('input'); fire('change');
      return 'ok';
    }
    el.focus();
    const proto = el.tagName === 'TEXTAREA' ? HTMLTextAreaElement.prototype : HTMLInputElement.prototype;
    const setter = Object.getOwnPropertyDescriptor(proto, 'value');
    if (setter && setter.set) setter.set.call(el, step.value || ''); else el.value = step.value || '';
    fire('input'); fire('change');
    return 'ok';
  }
  if (step.select !== undefined && step.select !== '') {
    if (el.tagName !== 'SELECT') return 'not_a_select';
    const want = (step.option || '').trim();
    const opt = [...el.options].find(o => o.textContent.trim() === want) || [...el.options].find(o => o.value === want);
    if (!opt) return 'option_unknown';
    el.value = opt.value;
    fire('input'); fire('change');
    return 'ok';
  }
  el.click();
  return 'ok';
})
