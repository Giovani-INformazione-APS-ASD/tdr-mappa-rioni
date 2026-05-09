# Torneo dei Rioni Address2Rione API

## Repository

La repository contiene il codice per il servizio di traduzione Indirizzo -> Rione di Carmiano o Magliano per il Torneo dei Rioni di Carmiano.

## API

L'API è esposta su `https://api.giovaniinformazione.it`.

## `/tdr-map/health`

Controlla se il servizio è attivo.
Esempio di risposta:

```json
{
  "message": "OK"
}
```

## `/tdr-map/get?address=<value>`

Fornito un indirizzo come `address` ritorna il Rione trovato insieme ad altre info accessorie.

Esempio di risposta con `address = Via Sara Librando 20`

```json
{
  "street": "Via Sara Librando",
  "number": "20",
  "latitude": "40.3321336",
  "longitude": "18.0493457",
  "rione": "San Giovanni"
}
```

In caso di errore, indisponibilità ritorna un errore.
Esempio con `address = Spero che potrate`:

```json
{
  "error": "Impossibile determinare Rione"
}
```

Il formato di errore è uguale per tutti gli errori tranne nel caso di `HTTP 429 Too Many Requests`
