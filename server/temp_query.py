import sqlite3

conn = sqlite3.connect(r'c:\Users\acfur\Downloads\nier\lunar-tear\server\db\game.db')
c = conn.cursor()
c.execute("SELECT sql FROM sqlite_master WHERE name IN ('user_materials', 'user_consumable_items');")
for r in c.fetchall():
    if r[0]: print(r[0])
c.execute("SELECT * FROM user_materials LIMIT 20;")
print("materials sample:", c.fetchall())

c.execute("SELECT * FROM user_consumable_items LIMIT 20;")
print("consumable items sample:", c.fetchall())
