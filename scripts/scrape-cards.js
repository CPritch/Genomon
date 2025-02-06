/* This scrapes the list of PTCGP cards and their effects from https://game8.co/games/Pokemon-TCG-Pocket/archives/482685 */

function parseAbility(node) {
  const abilitySpan = node.querySelector("td.left span.a-red");
  if (!abilitySpan) return;
  const title = abilitySpan.nextSibling.textContent.trim();
  let body = "";
  for (
    let x = abilitySpan.nextSibling.nextSibling;
    x.nodeName != "DIV";
    x = x.nextSibling
  ) {
    body += x.textContent.trim();
  }
  return {
    title: title,
    effect: body,
  };
}

function parseAttackCost(attackCostImages) {
  const costs = Array.from(attackCostImages).map((img) =>
    img.alt.replace("Pokemon TCG Pocket - ", "").trim()
  );
  return costs
    .map((cost) =>
      new Array(parseInt(cost.split(" ")[1] || 1)).fill(cost.split(" ")[0])
    )
    .flat();
}

function parseAttackEffect(attackElem) {
    let effect = "";
    let node = attackElem.nextSibling;
  
    while (node && !(node.classList && node.classList.contains("align"))) {
      if (node.nodeType === Node.TEXT_NODE) {
        let text = node.textContent.trim();
        if (text && !/^\d+$/.test(text)) {
          effect += (effect ? " " : "") + text;
        }
      }
      node = node.nextSibling;
    }
  
    return effect.trim() || null;
  }
  

function parsePokemonTable() {
  const rows = document.querySelectorAll(
    "div.archive-style-wrapper > div.scroll--table.table-header--fixed > table tr"
  );
  let cards = [];

  rows.forEach((row) => {
    const cells = row.querySelectorAll("td");
    if (cells.length < 10) return;
    const cardId = cells[1].textContent.trim().replaceAll(" ", "-");
    const cardName = cells[2].querySelector("a").textContent.trim();
    const rarity = cells[3].textContent.trim();
    const typeImg = cells[5].querySelector("img");
    const type = typeImg
      ? typeImg.alt.replace("Pokemon TCG Pocket - ", "").trim()
      : "";
    const stageText = cells[7].textContent.trim();
    const retreatCostImg = cells[9].querySelector(".align img");
    const retreatCost = retreatCostImg
      ? parseInt(retreatCostImg.getAttribute("width")) / 20
      : "";
    const ability = parseAbility(cells[9]);
    let requires = stageText != "Basic" ? cards.at(-1).id : null;
    if (
      cardName.indexOf("ex") > -1 ||
      cards.find((card) => card.name == cardName)
    ) {
      requires = cards.find(
        (card) => card.name == cardName.replace(/ ex$/, "")
      )?.requires;
    }
    let attacks = [];
    const attackElements = cells[9].querySelectorAll(".align");

    attackElements.forEach((attackElem) => {
      const attackNameElem = attackElem.querySelector("b");
      if (!attackNameElem) return;

      const attackName = attackNameElem.textContent.trim();
      if (attackName == "Retreat Cost") return;
      const attackCostImages = attackElem.querySelectorAll("img");
      const attackCost = parseAttackCost(attackCostImages);

      const attackDamageMatch = attackElem.nextSibling?.textContent
        ?.trim()
        .match(/\d+/);
      const attackDamage = attackDamageMatch
        ? parseInt(attackDamageMatch[0], 10)
        : 0;

      const attackEffect = parseAttackEffect(attackElem);

      attacks.push({
        name: attackName,
        cost: attackCost,
        damage: attackDamage,
        effect: attackEffect,
      });
    });
    if (type == "Supporter" || type == "Item" || type == "Pokemon Tool") {
      cards.push({
        id: cardId,
        name: cardName,
        rarity: rarity,
        type: type,
        ability: cells[9].textContent.trim()
      });
    } else {
      cards.push({
        id: cardId,
        name: cardName,
        rarity: rarity,
        type: type,
        stage: stageText,
        requires: requires,
        retreatCost: retreatCost,
        ability: ability,
        attacks: attacks,
      });
    }
  });

  return cards;
}

parsePokemonTable();
