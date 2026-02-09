const fs = require('fs');
const { Client } = require('pg');
const { parse } = require('csv-parse');

const connectionString = 'postgresql://artemz:dIZx9kCcDcugr06@109.71.244.236:5433/core_ui?sslmode=disable';

async function loadPrices() {
  const client = new Client({ connectionString });
  
  try {
    await client.connect();
    console.log('Подключено к PostgreSQL');

    const insertQuery = `
      INSERT INTO wb_product_prices (name, max_price, created_at, updated_at)
      VALUES ($1, $2, now(), now())
    `;

    const records = [];

    fs.createReadStream('prices.csv')
      .pipe(parse({
        delimiter: ';',
        columns: false,           // нет заголовка
        skip_empty_lines: true,
        trim: true,
      }))
      .on('data', (row) => {
        const [name, priceStr] = row;
        const price = parseInt(priceStr.trim(), 10);
        
        if (name && !isNaN(price)) {
          records.push(client.query(insertQuery, [name.trim(), price]));
        }
      })
      .on('end', async () => {
        try {
          await Promise.all(records);
          console.log(`Успешно загружено ${records.length} записей`);
        } catch (err) {
          console.error('Ошибка при вставке:', err);
        } finally {
          await client.end();
        }
      })
      .on('error', (err) => {
        console.error('Ошибка чтения CSV:', err);
        client.end();
      });

  } catch (err) {
    console.error('Ошибка подключения:', err);
    await client.end();
  }
}

loadPrices();
