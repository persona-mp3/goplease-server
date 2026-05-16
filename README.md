# 👾 go, please

trying to make a game in pure go

---
notes to self
### Phantom AP Rules

* **Base Pool:** Each living unit grants +1 AP to the team pool at the start of the round.
* **Phantom AP Generation:** If team counts are unequal, the disadvantaged team receives bonus points:
  `Phantom AP = Enemy Living Units - Friendly Living Units`
* **Individual Cap:** A single unit can expend a maximum of **2 AP per turn** (1 base + 1 phantom).
* **Dynamic Reset:** Phantom AP recalculates every round and drops to 0 as soon as team sizes equalize.

### Data Example

* **4v6 (Delta 2):** 4 Base + 2 Phantom = **6 AP**. All points can be spent (two units act twice).
* **2v6 (Delta 4):** 2 Base + 4 Phantom = 6 AP. Due to the 2 AP individual cap, the team can only spend **4 AP**. The remaining 2 Phantom AP waste. Sacrificing allies is inefficient.

---
