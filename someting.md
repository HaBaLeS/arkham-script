- Do Test (agility) level 3
    if fail get damage 1
    else nothing

---
- Do Test (brain) level 3
if fail get horror 1
else nothing

---
- Revelation Play card in threat
- once per round
    if (move, fight, evade)
    then action cost = 2

- end of round
    test (brain) level 3
    if success - remove card

--- 

- Revelation put card in threat area
- on play asset -> prevent play
- on play event -> prevent play
- end of Round Discard card

---
revelation place 1 doom
advance aganda if possible
---
on_revelation do_test(head, 4)
    success: noting
    fail: discard_asset_if_possible(1) or horror(2)

----

on_revelation:  
    play card to location (limit1)
        Location +2 shroud (nebel)

add event to location:
    on_investigate_success:
        discard card
---

# Events

on_revelation
on_play_card_type
on_end_of_round

# Actions
- Do test
- Place card in ThreatArea
- Remove card
- Play Doom Token
- AdvanceAgendaIfPossible