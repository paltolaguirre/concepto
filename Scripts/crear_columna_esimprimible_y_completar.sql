ALTER TABLE concepto ADD COLUMN esimprimible boolean;


UPDATE concepto SET esimprimible = true WHERE ID NOT IN (-21,-22,-23,-24,-25,-26,-27,-28)
