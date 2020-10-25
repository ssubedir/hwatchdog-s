CREATE DATABASE `hwatchdog` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

-- hwatchdog.hw_amazon_country definition

CREATE TABLE `hw_amazon_country` (
  `_id` int(11) NOT NULL AUTO_INCREMENT,
  `_country` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='location id of country';


-- hwatchdog.hw_cpu definition

CREATE TABLE `hw_cpu` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `AmazonId` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=88 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- hwatchdog.hw_cpu_amazon definition
CREATE TABLE `hw_cpu_amazon` (
  `_id` int(11) NOT NULL AUTO_INCREMENT,
  `_asin` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `_price` float DEFAULT NULL,
  `_country` int(11) DEFAULT NULL,
  `_time` datetime DEFAULT NULL,
  PRIMARY KEY (`_id`)
) ENGINE=InnoDB AUTO_INCREMENT=258 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='cpu prices amazon';
