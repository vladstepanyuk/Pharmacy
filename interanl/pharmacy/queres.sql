-- 1 ui y model y
-- 2 ui y model y
-- 3 ui y model y
-- 4 ui y model y
-- 5 ui y model y
-- 6 ui y model y
-- 7 ui y model y
-- 8 ui y model y
-- 9 ui y model y
-- 10 ui y model y
-- 12 ui y model y
-- 13 ui y model y

-- 1. Получить сведения о покупателях, которые не пришли забрать свой заказ
-- в назначенное им время и общее их число.

with cur_date as (select cast(current_date as DATE) date)
select distinct c.*
from consumer_order co
         cross join cur_date
         inner join recipe r on recipe_id = r.id
         inner join consumers c on consumer_id = c.id
where complete_time < cur_date.date
  and co.status = 'WAITING';


-- 2. Получить перечень и общее число покупателей, которые ждут прибытия
-- на склад нужных им медикаментов в целом и по указанной категории
-- медикаментов.

with need_types as (select id from drug_type)
select distinct c.*
from consumer_order co
         inner join recipe r on recipe_id = r.id
         inner join consumers c on consumer_id = c.id
         inner join drug d on r.drug_id = d.id
where co.status = 'WAITING-COMPS'
  and d.type_id in (select * from need_types);

-- 3. Получить перечень десяти наиболее часто используемых медикаментов в
-- целом и указанной категории медикаментов.

with need_types as (select distinct *
                    from (VALUES (1), (2), (3))),
     help_table as (select comp_id, count(tech_id)
                    from technology_components tc
                             inner join drug d on tc.comp_id = d.id
                    where d.type_id in (select * from need_types)
                    group by comp_id)
select d.id, d.name, d.type_id, d.critical_count, d.price, coalesce(ht.count, 0) count
from drug d
    left join help_table ht
on ht.comp_id = d.id
where d.type_id in (select * from need_types)
order by count DESC
    LIMIT 10;



-- 4. Получить какой объем указанных веществ использован за указанный
-- период.

with need_drug_ids as (select distinct *
                       from (VALUES (1), (2), (3))),
     need_drugs as (select id
                    from drug d
                    where id in (select * from need_drug_ids)),
     period as (select date ('2023-12-20') "start", date ('2023-12-25') "end"), drug_count as (
select drug_id, count (count) count
from inventasrization
    cross join period p
where p.start <= date
  and date <= p."end"
  and drug_id in (select * from need_drug_ids)
group by drug_id)
select id, coalesce(dc.count, 0) count
from need_drugs
    left join drug_count dc
on id = dc.drug_id;


-- 5. Получить перечень и общее число покупателей, заказывавших
-- определенное лекарство или определенные типы лекарств за данный
-- период.

with need_drug_ids as (select distinct *
                       from (VALUES (1), (2), (3))),
     need_types as (select distinct *
                    from (VALUES (1), (3))),
     period as (select date ('2023-12-20') "start", date ('2023-12-25') "end")
select c.id, c.name
from consumer_order co
         cross join period
         inner join recipe r on co.recipe_id = r.id
         inner join consumers c on c.id = r.consumer_id
         inner join drug d on r.drug_id = d.id
where (drug_id in (select * from need_drug_ids) or d.type_id in (select * from need_types))
  and period.start <= co.order_date
  and co.order_date <= period."end";



-- 6. Получить перечень и типы лекарств, достигших своей критической нормы
-- или закончившихся.

select id, type_id, name
from drug
         left join storage s on drug.id = s.comp_id
where coalesce(count, 0) <= drug.critical_count;


-- 7. Получить перечень лекарств с минимальным запасом на складе в целом и
-- по указанной категории медикаментов.

with need_types as (select distinct *
                    from (VALUES (1), (2))),
     drug_count as (select id, coalesce(count, 0) count
from drug
    left join storage s
on drug.id = s.comp_id
where coalesce (count
    , 0) <= drug.critical_count
  and id in (select * from need_types))
select d.id, name
from drug_count
         inner join drug d on drug_count.id = d.id
where count = (select min(count) from drug_count);


-- 8. Получить полный перечень и общее число заказов находящихся в
-- производстве.

select co.recipe_id, co.id, co.complete_time
from consumer_order co
         inner join recipe r on recipe_id = r.id
where co.status = 'IN-PROGRESS';

-- 9. Получить полный перечень и общее число препаратов требующихся для
-- заказов, находящихся в производстве.


with drugs_comp as (select td.drug_id, tc.comp_id, tc.count
                    from technology_components tc
                             inner join technology_drug td on tc.tech_id = td.tech_id)
select distinct dc.comp_id
from consumer_order co
         inner join recipe r on r.id = co.recipe_id
         inner join drugs_comp dc on dc.drug_id = r.drug_id
         inner join drug d on dc.drug_id = d.id
where co.status = 'IN-PROGRESS';


-- 10. Получить все технологии приготовления лекарств указанных типов,
-- конкретных лекарств, лекарств, находящихся в справочнике заказов в
-- производстве.


with drugs_in_progress as (select distinct r.drug_id
                           from consumer_order co
                                    inner join recipe r on r.id = co.recipe_id
                           where co.status = 'IN-PROGRESS'),
     need_drug_ids as (select distinct *
                       from (VALUES (1), (2), (3))),
     need_types as (select distinct *
                    from (VALUES (1), (3)))
select tb.id, tb.description
from technology_drug td
         inner join technology_book tb on td.tech_id = tb.id
where td.drug_id in (select * from drugs_in_progress)
   or td.drug_id in (select * from need_types)
   or td.drug_id in (select * from need_drug_ids);


-- 11. Получить сведения о ценах на указанное лекарство в готовом виде, об
-- объеме и ценах на все компоненты, требующиеся для этого лекарства.

with need_drug_ids as (select distinct *
                       from (VALUES (1), (2), (3)))
select id, name, price
from drug;


with need_drug_ids as (select distinct *
                       from (VALUES (1), (2), (3))),
     id_sum_stat as (select tc.comp_id, sum(tc.count) sum_count, sum(tc.count * d.price) sum_price
                     from technology_components tc
                              inner join technology_drug td on tc.tech_id = td.tech_id
                              inner join drug d on tc.comp_id = d.id
                     where drug_id in (select * from need_drug_ids)
                     group by tc.comp_id)
select name, sum_count, sum_price
from id_sum_stat
         inner join drug d on comp_id = d.id;

-- 12. Получить сведения о наиболее часто делающих заказы клиентах на
-- медикаменты определенного типа, на конкретные медикаменты.

with need_drug_ids as (select distinct *
                       from (VALUES (1), (2), (3))),
     need_types as (select distinct *
                    from (VALUES (1), (3))),
     help_table_1 as (select r.consumer_id, count(*) count
from consumer_order co
    inner join recipe r
on co.recipe_id = r.id
    inner join drug d on r.drug_id = d.id
where d.id in (select * from need_drug_ids)
   or d.type_id in (select * from need_types)
group by r.consumer_id),
    consumer_count as (
select id, coalesce (ht.count, 0) count
from consumers
    left join help_table_1 ht
on id = ht.consumer_id)
select c.id, c.name, cc.count
from consumer_count cc
         inner join consumers c on c.id = cc.id
where count = (select max(count) from consumer_count);

-- 13. Получить сведения о конкретном лекарстве (его тип, способ
-- приготовления, названия всех компонент, цены, его количество на
-- складе).


with need_drug_ids as (select int (6) "id")
select d.name_, d.type_id, d.critical_count, d.price, coalesce(tb.description, 'nil')
from drug d
         inner join technology_drug td on d.id = td.drug_id
         inner join technology_book tb on td.tech_id = tb.id
where d.id in (select * from need_drug_ids);

with need_drug_ids as (select 2)
select d.id, d.name_, count, count * price price
from technology_drug td
         inner join technology_components tc on td.tech_id = tc.tech_id
         inner join drug d on tc.comp_id = drug.id
where td.drug_id in (select * from need_drug_ids);
