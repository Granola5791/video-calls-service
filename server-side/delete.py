x = "\"\"%' delete from products where name like '%"
query = "select * from products where name like '%" + x + "%'" 
print(query)