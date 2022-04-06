// Package promqlext расширяет язык PromQL новыми функциями
//
// Часть функций перенесена для обратной совместимости с OkQL, языком запросов системы okmeter.
// Описание исходных функций доступна в документации проекта docs: https://okmeter.io/misc/docs#query-language,
// а реализация — https://fox.flant.com/okmeter/okmeter/-/blob/master/okmeter/core/eval_funcs.py.
//
// Функции регистрируются в глобальных переменных пакета promql с помощью специальных регистраторов.
package promqlext
