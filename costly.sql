CREATE TABLE IF NOT EXISTS public.items (
  id         BIGSERIAL PRIMARY KEY,
  name       TEXT NOT NULL,                     -- user-facing label
  price      NUMERIC(12,2) NOT NULL,            -- unit price (LKR)
  quantity   INTEGER NOT NULL,                  -- units in this row
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  -- canonicalized name used only for dedup / uniqueness
  name_lc    TEXT GENERATED ALWAYS AS (LOWER(TRIM(name))) STORED,

  -- convenience; not required for aggregates, but nice for reads
  line_total NUMERIC(14,2)
    GENERATED ALWAYS AS (price * quantity) STORED,

  -- data guards
  CONSTRAINT price_non_negative   CHECK (price    >= 0),
  CONSTRAINT quantity_positive    CHECK (quantity >= 1)
);

-- one logical row per (canonical name, price)
CREATE UNIQUE INDEX IF NOT EXISTS uniq_items_name_lc_price
  ON public.items (name_lc, price);
